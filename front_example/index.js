const express       = require('express');
const app           = express();
const path          = require('path');
const bodyParser    = require('body-parser');
const cookieParser  = require('cookie-parser');
const nunjucks      = require('nunjucks');
const cors 			= require('cors')
const logger 		= require('pino')();
const request		= require('got');
const routes        = require('./routes');
const ApiBridge     = require('./apibridge');

// Ignore self-signed certificates
process.env["NODE_TLS_REJECT_UNAUTHORIZED"] = 0;

// Set environment variables
host = process.env.HOST_API || 'localhost';
port = process.env.PORT || 8080;

// Create a new instance of the API Bridge
const ab = new ApiBridge(request, `https://${host}:5000`, logger);

// Configuring CORS (Cross-Origin Resource Sharing) middleware
const corsOptions = {
	origin: '*',
	optionsSuccessStatus: 200,
};
app.use(cors(corsOptions))

// Body parser setup (to parse JSON and URL-encoded request bodies)
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(cookieParser());

// Nunjucks setup (for rendering HTML templates)
nunjucks.configure('views', {
	autoescape: true,
	express: app
});

// Static files setup
app.set('views', './views');
app.use('/static', express.static(path.resolve('static')));

// Routes setup
app.use(routes(ab, corsOptions));

// Error handling middleware
app.all('*', (req, res) => {
	return res.status(404).send({
		message: '404 page not found'
	});
});

(async () => {
	logger.info(`Server started on port ${port}`);
	app.listen(port, '0.0.0.0', () => console.log(`Listening on port ${port}`));
})();