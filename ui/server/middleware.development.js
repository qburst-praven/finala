const path = require("path");
const webpack = require("webpack");
const webpackDevMiddleware = require("webpack-dev-middleware");
const webpackHotMiddleware = require("webpack-hot-middleware");
const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function middleware(app, config) {
  const compiler = webpack(config);

  // Route for frontend to get API endpoint configuration (BEFORE general /api proxy)
  app.get("/api/v1/settings", (req, res) => {
    // Default to localhost:8089 if env var is not set by the Go finala ui process
    const apiEndpoint = process.env.API_ENDPOINT || "http://localhost:8089"; 
    console.log(`[SETTINGS] Serving API endpoint: ${apiEndpoint}`);
    res.send({ api_endpoint: apiEndpoint });
  });

  // Proxy API requests (handles all other /api/... calls)
  app.use(
    '/api', // This will not catch /api/v1/settings if it's defined above
    createProxyMiddleware({
      target: 'http://localhost:8089',
      changeOrigin: true,
      logLevel: 'debug',
      onProxyReq: function(proxyReq, req, res) {
        console.log(`[HPM] PROXY: ${req.method} ${req.url} -> ${proxyReq.getHeader('host')}${proxyReq.path}`);
      },
      onError: function(err, req, res) {
        console.error('[HPM] PROXY ERROR:', err);
        if (!res.headersSent) {
          res.writeHead(500, { 'Content-Type': 'text/plain' });
        }
        res.end('Proxy error.');
      },
    })
  );

  // THEN, Webpack dev middleware
  const devMiddlewareInstance = webpackDevMiddleware(compiler, {
    publicPath: config.output.publicPath,
    historyApiFallback: true,
  });
  app.use(devMiddlewareInstance);

  // THEN, Webpack hot middleware
  app.use(webpackHotMiddleware(compiler));

  const fs = devMiddlewareInstance.fileSystem;

  // Catch-all for SPA routing, ensuring it doesn't conflict with /api proxy or /api/v1/settings
  app.get("*", (req, res, next) => {
    if (req.path.startsWith('/api/')) { // Do not interfere with API calls
      return next();
    }
    fs.readFile(path.join(compiler.outputPath, "index.html"), (err, file) => {
      if (err) {
        res.sendStatus(404);
      } else {
        res.send(file.toString());
      }
    });
  });
};
