package restful

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"
	"ragbox/common"
	. "ragbox/config"
	"ragbox/embedding"
	"ragbox/llm"
	"ragbox/reader"
	"ragbox/rerank"
	. "ragbox/store"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type RestfulServer struct {
	engine *gin.Engine
	addr   string
}

func (srv *RestfulServer) Run() {
	srv.engine.Run(srv.addr)
}

type Document struct {
	ID              int    `json:"id"`
	OriginFilename  string `json:"origin_filename"`
	MimeType        string `json:"mime_type"`
	CreateTimestamp string `json:"create_timestamp"`
	UpdateTimestamp string `json:"update_timestamp"`
}

type ChatRequest struct {
	Content     string `json:"content"`
	EnableKbase bool   `json:"enableKbase"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	Passwd          string `json:"password"`
	CreateTimestamp int64  `json:"create_timestamp"`
	UpdateTimestamp int64  `json:"update_timestamp"`
}

func getFileLst(ctx context.Context, pageNo, pageSize int) ([]*Document, error) {
	rows, err := Store.Query(ctx, "select id, origin_filename, mimetype, create_timestamp, update_timestamp from documents order by update_timestamp desc, id desc limit ?, ?", pageNo*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Document
	for rows.Next() {
		doc := new(Document)
		if err := rows.Scan(&doc.ID, &doc.OriginFilename, &doc.MimeType, &doc.CreateTimestamp, &doc.UpdateTimestamp); err != nil {
			return nil, err
		}
		result = append(result, doc)
	}

	return result, rows.Err()
}

func LstFile(c *gin.Context) {
	flist, err := getFileLst(c, 0, 10)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": flist,
		"msg":  "success",
	})
}

func queryDocumentIdByFileId(ctx context.Context, fileId string) (string, error) {
	row, err := Store.Query(ctx, "select doc_id from documents where id = ?", fileId)
	if err != nil {
		return "", err
	}
	defer row.Close()

	if row.Next() {
		var docId string
		if err := row.Scan(&docId); err != nil {
			return "", err
		}
		return docId, nil
	}
	return "", errors.New("document not found")
}

func querySegmentsByDocId(ctx context.Context, docId string) ([]string, error) {
	return VectorStore.QueryContentsByDocID(ctx, "knowledgebase", docId)
}

func GetSegementLst(c *gin.Context) {
	fileid := c.Param("fileid")
	docId, err := queryDocumentIdByFileId(c, fileid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	result, err := querySegmentsByDocId(c, docId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": result,
		"msg":  "success",
	})
}

func genVector(section string) ([]float32, error) {
	var model embedding.EmbeddingModel
	model = new(embedding.BailianEmbedding)
	return model.GetEmbedding(section)
}

func recall(query string) ([]string, error) {
	vector, err := genVector(query)
	if err != nil {
		return nil, err
	}
	slog.Debug("recall", "vector", vector)

	results, err := VectorStore.Search(context.Background(), "knowledgebase", vector, 5)
	if err != nil {
		return nil, err
	}

	var segments []string

	for _, result := range results {
		col_content := result.GetColumn("content")
		for i := 0; i < result.ResultCount; i++ {
			content, err := col_content.GetAsString(i)
			if err != nil {
				return nil, err
			}
			segments = append(segments, content)
		}
	}

	return segments, nil
}

func callRerank(query string, segments []string) ([]string, error) {
	var model rerank.RerankModel
	model = new(rerank.BailianRerank)
	return model.Rerank(query, segments)
}

func Chat(c *gin.Context) {
	// content := "抱歉，我无法直接访问您的阿里云控制台或数据库实例 `rm-3yd2hd567avlw1prb` 的内部巡检数据，因为我没有您的云账号权限，也无法实时获取该实例的运行状态。\n\n不过，我可以为您提供**关于如何查看该实例巡检结果**以及**如何评估巡检结果**的标准方法。您可以根据以下指引自行操作：\n\n### 1. 如何获取该实例的巡检结果？\n您需要登录**阿里云RDS管理控制台**，按以下路径查找：\n\n- 登录 [阿里云控制台](https://rds.console.aliyun.com/)。\n- 在左侧导航栏选择“实例列表”，找到 `rm-3yd2hd567avlw1prb` 并点击进入实例详情页。\n- 查看 **“健康巡检”** 或 **“自治中心”**（通常在“数据库自治服务DAS”或“服务可用性”菜单下）。\n- 如果使用的是**数据库自治服务（DAS）**，您可以看到“一键巡检”或“智能巡检”报告，系统会根据 CPU、内存、磁盘、连接数、慢SQL等指标生成健康评分。\n\n### 2. 巡检结果评估的关键维度\n拿到巡检报告后，请重点关注以下核心指标，并对照阈值进行评估：\n\n| 评估维度 | 健康标准（建议值） | 风险提示（需处理） |\n| :--- | :--- | :--- |\n| **CPU使用率** | 平均值 < 70%，峰值 < 90% | 持续 > 90% 可能导致SQL超时或锁等待 |\n| **内存使用率** | < 80% | 内存溢出（OOM）可能导致实例重启 |\n| **磁盘空间** | 使用率 < 80%，且剩余空间 > 5GB | 使用率 > 90% 会导致实例锁定为只读 |\n| **连接数** | 使用率 < 80% | 接近 max_connections 上限会导致应用报错 |\n| **慢SQL（慢查询）** | 无大量长事务或慢查询积压 | 单条SQL执行时间 > 1秒且频繁出现，需优化索引 |\n| **主备延迟**（高可用） | 延迟 < 1秒 | 延迟持续高可能导致灾备切换丢数据 |\n| **备份状态** | 最近一天内备份成功 | 备份失败且无日志，存在数据丢失风险 |\n| **死锁/锁等待** | 无死锁或锁等待超时 | 业务并发高时出现死锁，需排查应用逻辑 |\n\n### 3. 如果发现异常，建议的应对措施\n- **评分低于 80分**：说明实例存在资源瓶颈，建议查看具体失分项（如慢SQL、磁盘扩容）。\n- **有告警事件**：在“事件管理”中查看具体告警（如“实例可用性下降”、“磁盘水位告警”），按提示处理。\n- **如果巡检提示“实例规格不足”**：建议根据业务高峰期评估是否升级CPU/内存，或开启**自动弹性伸缩**（如Serverless）。\n- **如果有慢SQL**：点击慢SQL明细，使用“SQL洞察”优化对应语句，或添加合适的索引。\n\n---\n\n**重要提醒**：由于我无法访问您的具体数据，**如果巡检结果中出现了“严重”级别（红码）的告警，请您立即登录控制台确认，或联系阿里云技术支持（提交工单）**，以免影响业务。\n\n如果您已经拿到了具体的巡检报告，**可以把报告中的“异常项”文字内容（隐去敏感IP/账号）发给我**，我可以帮您分析具体风险和优化建议。"
	req := new(ChatRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	payload := "用户问题: " + req.Content
	if req.EnableKbase {
		payload += "\n\n请根据知识库内容回答用户问题，不要编造答案。"

		recall, err := recall(req.Content)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		reranked, err := callRerank(req.Content, recall)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(reranked) > 0 {
			payload += "\n\n知识库内容:\n"
			for i, segment := range reranked {
				payload += "段落" + strconv.Itoa(i+1) + ": " + segment + "\n"
			}
		} else {
			payload += "\n\n知识库中没有相关内容。"
		}

	}

	// response := llm.LLMResponse{}
	response, err := llm.Send2LLM(payload)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(response.Choices) == 0 {
		c.JSON(500, gin.H{
			"error": "LLM returned empty choices",
		})
		return
	}
	c.JSON(200, gin.H{
		"data": response.Choices[0].Message.Content,
		// "data": content,
		"msg": "success",
	})
}

func fileHandler(c *gin.Context) {

	file, _ := c.FormFile("file")
	f, _ := file.Open()
	defer f.Close()

	// 读取前 512 字节
	buf := make([]byte, 512)
	n, _ := f.Read(buf)

	mimeType := http.DetectContentType(buf[:n])

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if mimeType == "application/zip" && ext == ".docx" {
		mimeType = "##word##"
	}

	slog.Info("file upload.", "mimeType", mimeType)

	documentId := uuid.New().String()

	storagePath := Config.Storage.File.Path + "/"
	dst := storagePath + documentId
	c.SaveUploadedFile(file, dst)

	if _, err := Store.Insert(c, "insert into documents(filepath, origin_filename, mimetype, doc_id, user_id, create_timestamp, update_timestamp) value(?,?,?,?,?,?,?);", dst, file.Filename, mimeType, documentId, 2, time.Now().Unix(), time.Now().Unix()); err != nil {
		slog.Error("insert document error.", "err", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	go func() {
		reader, err := reader.NewInspctRptReader(dst)
		if err != nil {
			slog.Error("new inspection reader error.", "err", err)
			return
		}
		sections, err := reader.Read()
		if err != nil {
			slog.Error("read inspection report error.", "err", err)
			return
		}

		slog.Info("The file sharding is complete.", "documentId", documentId)

		var model embedding.EmbeddingModel
		model = new(embedding.BailianEmbedding)

		for _, section := range sections {
			slog.Debug("section", "section", section)
			vector, err := model.GetEmbedding(section)
			if err != nil {
				panic(err)
			}
			if len(vector) == 0 {
				panic("embedding is empty")
			}
			slog.Debug("embedding", "embedding", vector)

			if err := VectorStore.Insert(c, "knowledgebase", []map[string]any{{
				"document_id": documentId,
				"content":     section,
				"embedding":   vector,
			}}); err != nil {
				slog.Error("error on insert to vector database", "error", err)
			}
		}
	}()

	c.JSON(200, gin.H{"message": "upload success", "data": map[string]any{"documentId": documentId, "fname": file.Filename, "mimeType": mimeType, "length": file.Size}})
}

// HashPassword 使用 bcrypt 对密码进行加盐哈希
func HashPassword(password string) (string, error) {
	// GenerateFromPassword 自动生成盐值并进行哈希
	// bcrypt.DefaultCost 是默认的计算权重（当前为 10），可根据服务器性能调整
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证用户输入的明文密码与数据库中的哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func FindUserWithUsername(c context.Context, username string) (*User, error) {
	rows, err := Store.Query(c, "select id, username, passwd, create_timestamp, update_timestamp from users where username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Passwd, &user.CreateTimestamp, &user.UpdateTimestamp); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

func genToken() string {
	return uuid.New().String()
}

func loginSuccessHandler(user *User) (string, error) {
	token := genToken()
	value, err := common.JsonStringify(user)
	if err != nil {
		slog.Error("error on json stringify user", "error", err)
		return "", err
	}
	GlobalCache.SetWithExpire(token, value, time.Minute*30) // 设置 token 过期时间为 30 分钟
	return token, nil
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := FindUserWithUsername(c, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user == nil || !CheckPasswordHash(req.Password, user.Passwd) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	// 处理登录逻辑，例如验证用户名和密码
	token, err := loginSuccessHandler(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "success",
		"token": token,
	})
}

func registerRouter(engine *gin.Engine) {
	apiV1 := engine.Group("/api_v1")
	apiV1.GET("/files", LstFile)
	apiV1.POST("/chat", Chat)
	apiV1.GET("/file/:fileid/segments", GetSegementLst)
	apiV1.POST("/upload", fileHandler)
	apiV1.POST("/login", login)
}

func NewRestfulServer(addr string) *RestfulServer {

	engine := gin.Default()
	registerRouter(engine)

	return &RestfulServer{engine: engine, addr: addr}

}
