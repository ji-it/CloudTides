import { create, router, defaults, bodyParser, rewriter } from 'json-server';
import { runInterceptors } from './interceptors'

const dbRouter = router('db-tmp.json');
const customRoutes = require('./routes.json');

// init
const server = create();
server.use(defaults());
server.use(bodyParser);

// Router
server.use(rewriter(customRoutes));
runInterceptors(server);
server.use(dbRouter);

server.listen(3000, () => {
  console.log('JSON Server is running')
});
