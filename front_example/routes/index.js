const path           = require('path');
const express        = require('express');
const router         = express.Router();
const cors 			 = require('cors');
const { Session } = require('inspector');


let ab;
let crs;

const response = data => ({ message: data });

router.get('/', cors(crs), (req, res, next) => {
	return res.render('index.html');
});

router.post('/api/register', cors(crs), async (req, res, next) => {
	const { username, password } = req.body;
	if (username && password) {
		return ab.register(username, password)
		.then(()  => res.send(response('User registered successfully!')))
		.catch(() => res.status(401).send(response('Something went wrong!')));
	}
	return res.status(401).send(response('Please fill out all the required fields!'));
});

router.post('/api/login', cors(crs), async (req, res, next) => {
	const { username, password } = req.body;
	if (username && password) {
		return ab.login(username, password)
		.then(token => {
			res.cookie('session', token, { maxAge: 3600000 });
			return res.send(response('User authenticated successfully!'));
		})
		.catch(() => res.status(403).send(response('Invalid username or password!')));
	}
	return res.status(500).send(response('Missing parameters!'));
});

router.post('/api/submit', cors(crs), async (req, res, next) => {
	const { message, chatDesc } = req.body;
	if (message && chatDesc) {
        return ab.sendMessage(req.cookies.session, chatDesc, message)
        .then(() => res.send(response('Message sent successfully!')))
        .catch(() => res.status(500).send(response('Internal server error')));
    }
	return next();
})

router.post('/api/chat/create/:chatDesc', cors(crs), async (req, res, next) => {
	const { chatDesc } = req.params;
    if (chatDesc) {
        return ab.createChat(req.cookies.session, chatDesc)
       .then(() => res.send(response('Chat created successfully!')))
       .catch(() => res.status(500).send(response('Internal server error')));
    }
    return res.status(400).send(response('Missing chat description!'));
});

router.post('/api/chat/follow/:chatDesc', cors(crs), async (req, res, next) => {
	const { chatDesc } = req.params;
    if (chatDesc) {
        return ab.followChat(req.cookies.session, chatDesc)
       .then(() => res.send(response('Chat joined successfully!')))
       .catch(() => res.status(500).send(response('Internal server error')));
    }
    return res.status(400).send(response('Missing chat description!'));
});

router.post('/api/upload', cors(crs), async (req, res, next) => {
	if (!req.files || Object.keys(req.files).length === 0) {
        return res.status(400).send(response('No files were uploaded.'));
    }
    const file = req.files.file;
    if (!file.mimetype.startsWith('image/')) {
        return res.status(400).send(response('Only image files are allowed.'));
    }
    file.mv(path.join(__dirname, '../../public/uploads/' + file.name), err => {
        if (err) return res.status(500).send(response('Error uploading file.'));
        return res.send(response('File uploaded successfully.'));
    });
	
})

router.get('/api/avatar/:username', cors(crs), (req, res, next) => {
	if (req.cookies.session === undefined) return res.redirect('/');
    return ab.getAvatar(req.cookies.session, req.params.username)
    .then(avatar => { res.sendFile(path.join(__dirname, '../../public/uploads/' + avatar)); })
    .catch(() => { res.status(404).send('Avatar not found.'); });
})

router.get('/chat/:chatDesc', cors(crs), (req, res, next) => {
	if (req.cookies.session === undefined) return res.redirect('/');
    return ab.getMessages(req.cookies.session, req.params.chatDesc, req.query.offset, req.query.limit)
    .then(messages => { res.render('messenger.html', { messages }); })
    .catch(() => { res.redirect('/'); });
})

router.get('/chats', (req, res) => {
	if (req.cookies.session === undefined) return res.redirect('/');
	return ab.getInfoChats(req.cookies.session, req.query.offset, req.query.limit)
	.then(chats => {
		res.render('chats.html', { chats });
	})
	.catch(() => { res.redirect('/'); });
});

router.get('/profile', async (req, res, next) => {
	return ab.getInfoSelf(req.cookies.session)
	.then(user => {
			if(user === undefined) return res.redirect('/');
			res.render('profile.html', { user });
	})
	.catch(() => res.redirect('/'));
});

router.get('/logout', (req, res) => {
	ab.logout(res.cookie('session'));
	res.clearCookie('session');
	return res.redirect('/');
});

module.exports = (apib, cors) => {
	ab = apib
	crs = cors
	return router;
};
