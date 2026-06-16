import axios from 'axios'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

http.interceptors.response.use(
  (res) => {
    if (res.data.code !== 0) {
      return Promise.reject(new Error(res.data.message || '请求失败'))
    }
    return res.data
  },
  (err) => Promise.reject(err),
)

export default http
