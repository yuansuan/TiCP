import axios from 'axios'

export function exportFile(params, name) {
  axios({
    url: '/api/v1/auditlog/exportAll',
    method: 'POST',
    data: params,
    responseType: 'blob',
  }).then(response => {
    const url = window.URL.createObjectURL(new Blob([response.data]))
    // 取后端给前端返的请求头中的文件名称
    const temp =response.headers["content-disposition"].split(";")[1].split("filename=")[1];
    const fileName = decodeURIComponent(temp)

    const temps = fileName.split('.')
    temps.shift()

    const link = document.createElement('a')
    link.href = url
    link.setAttribute(
      'download',
      name ? `${name}.${temps.join('.')}` : fileName
    )
    document.body.appendChild(link)
    link.click()
  })
}
