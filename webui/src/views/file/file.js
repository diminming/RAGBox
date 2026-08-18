import request from '@/utils/request'

export class File {
    constructor(fileId) {
        this.id = fileId
    }

    static getLstByPage(pageNo, pageSize) {
        return new Promise((resolve, reject) => {
            request.get(`/files?pageNo=${pageNo}&pageSize=${pageSize}`)
                .then(data => {
                    resolve(data);
                })
                .catch(error => {
                    reject(error);
                });
        });
    }

    lstSegments() {
        return new Promise((resolve, reject) => {
            request.get(`/file/${this.id}/segments`)
                .then(data => {
                    resolve(data);
                })
                .catch(error => {
                    reject(error);
                });
        });
    }
}
