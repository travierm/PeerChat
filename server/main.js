'use strict';
var Cache = require('./cache');

const Hapi = require('hapi');
const server = Hapi.Server({
    port: 3000,
    host: 'localhost'
});

const signalCache = Cache();
const answerCache = Cache();

server.route({
    method: 'GET',
    path: '/',
    handler: () => {
        return "<pre>" + JSON.stringify(signalCache.data) + "</pre>"
    }
});

server.route({
    method: 'GET',
    path: '/signal/{id}',
    config: {
        cors: {
            origin: ['*'],
            additionalHeaders: ['cache-control', 'x-requested-with']
        }
    },
    handler: (request) => {
        const signal = signalCache.get(request.params.id);

        if(!signal) {
            return "bad id";
        }

        return JSON.stringify(signal);
    }
});

server.route({
    method: 'POST',
    path: '/signal',
    config: {
        cors: {
            origin: ['*'],
            additionalHeaders: ['cache-control', 'x-requested-with']
        }
    },
    handler: (request) => {
        const { id, signal } = request.payload;

        console.log("Cached new signal #" + id);

        signalCache.set(id, signal);

        return "success";
    }
})

server.route({
    method: 'POST',
    path: '/answer',
    config: {
        cors: {
            origin: ['*'],
            additionalHeaders: ['cache-control', 'x-requested-with']
        }
    },
    handler: (request) => {
        const { id, signal } = request.payload;

        console.log("Cached new answer #" + id);

        answerCache.set(id, signal);

        return "success";
    }
})

server.route({
    method: 'GET',
    path: '/answer/{id}',
    config: {
        cors: {
            origin: ['*'],
            additionalHeaders: ['cache-control', 'x-requested-with']
        }
    },
    handler: (request) => {
        const { id } = request.params;
        
        var data = answerCache.get(id);

        if(!data) {
            return "bad id";
        }

        return JSON.stringify(data);
    }
})

const init = async () => {

    await server.start();
    console.log(`Server running at: ${server.info.uri}`);
};

process.on('unhandledRejection', (err) => {

    console.log(err);
    process.exit(1);
});

init();
