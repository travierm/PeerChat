import axios from 'axios';
import nanoid from 'nanoid';

var Peer = require('simple-peer')

var p = new Peer({ initiator: location.hash === '#create', trickle: false })



p.on('signal', function (data) {

  axios.post('http://localhost:3000/signal', {
      id: nanoid(10),
      signal: JSON.stringify(data)
  });


})

p.on('connect', function () {
  console.log('CONNECT')
  p.send('whatever' + Math.random())
})

p.on('data', function (data) {
  console.log('data: ' + data)
})


export default function sendSignal(data) {
    p.signal(JSON.parse(data))
}