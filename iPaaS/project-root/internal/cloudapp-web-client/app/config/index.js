const devConfig = {
  port: 8001,
  retry_limit: 10,
  retry_interval_millsec: 1000,
  clipboard_buf_sz: 1000,
  input_ws: {
    jsonrpc_method: 'INPUT',
    sub_url: 'input',
  },
  webrtc: {
    ice_servers: [
      { urls: 'stun:stun.services.mozilla.com' },
      { urls: 'stun:stun.l.google.com:19302' },
    ],
  },
}

const prodConfig = {
  ...devConfig,
  // port:8000 is used for demo environment
  port: 8000,
  webrtc: {
    ice_servers: [
      {
        urls: ['turn:115.159.149.167:3478?transport=udp'],
        username: 'root',
        credential: 'lambdacal',
      },
    ],
  },
}

const getConfig = () => {
  if (G_ENV === 'dev') {
    return devConfig
  }
  return prodConfig
}

export default getConfig()
