// 使用环境变量判断是否为开发环境，在构建生产环境代码时，相应代码块会被删掉
const isDev = process.env.NODE_ENV === 'development'
export default isDev
