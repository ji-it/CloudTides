// @ts-check


const {
  npm_package_config_port: PORT = 4200,
  start_mode: MODE = 'local',
} = process.env;

const endpoint = {
  local: 'http://localhost:3000',
  e2e: 'http://localhost:8080',
  test: 'http://10.163.95.185:30220',
};

const target = endpoint[MODE];

console.log(`
==========================================
Starting Mode: [${MODE}]

  /v1/       -> ${target}
  /api-local/ -> ${endpoint.local}
==========================================
`);

const PROXY_CONFIG = {

  '/api-local/': {
    target: endpoint.local,
    changeOrigin: true,
    secure: false,
    pathRewrite: {
      '^/api-local': '/api',
    },
  },

  '/v1/': {
    target: target,
    changeOrigin: true,
    secure: false,
    pathRewrite: {
      '^/v1': '/v1',
    },
    onProxyRes(proxyResponse) {
      const cookies = proxyResponse.headers['set-cookie'];
      const prune = (cookie = '') => cookie.replace(/;\W*secure/gi, '');

      if (cookies) {
        proxyResponse.headers['set-cookie'] = cookies.map(prune);
      }
    }
  },

}

module.exports = PROXY_CONFIG;
