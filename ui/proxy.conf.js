// @ts-check


const {
  npm_package_config_port: PORT = 4200,
  start_mode: MODE = 'local',
} = process.env;

const endpoint = {
  local: 'http://localhost:3000',
  e2e: 'http://localhost:8080'
};

const target = endpoint[MODE];

console.log(`
==========================================
Starting Mode: [${MODE}]

  /api/       -> ${target}
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

  '/api/': {
    target: target,
    changeOrigin: true,
    secure: false,
    pathRewrite: {
      '^/api': '/api',
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
