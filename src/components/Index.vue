<template>
  <body>
      <a href="#create" onclick="location.reload()">Create Session</a>

      <p>Join Session [{{inputHash}}]</p>
      <input v-model="inputHash" /><button @click="joinSession">Join</button>

      <p>Messages:</p>

      <pre id="outgoing">
        <p v-for="msg in log">{{msg}}</p>
      </pre>

      <input v-model="inputMessage" />
      <button @click="send">Send</button>
  </body>
</template>

<script>
import Peer from 'simple-peer';
import nanoid from 'nanoid';
import axios from 'axios';

export default {
  name: 'Index',
  props: {
  },
  data() {
      return {
          peer: false,
          blockLoop: false,
          inputHash: "",
          inputMessage: "",
          log: [
            "Welcome to PeerChat"
          ]
      }
  },
  mounted() {
    this.peer = initPeer(this);
  },
  methods: {
    send() {
      this.log.push(this.inputMessage);
      this.peer.send(this.inputMessage);
      this.inputMessage = "";
    },
    startAnswerCheck() {
      let self = this;
      setInterval(function() {
        self.answerCheck();
      }, 5000);
    },
    answerCheck() {
      if(this.blockLoop) {
        return;
      }

      axios.get('http://localhost:3000/answer/' + this.inputHash)
        .then((resp) => {
          if(resp.data !== "bad id") {

            console.log("DOING 2nd Signal");
            this.peer.signal(JSON.parse(resp.data));
            this.blockLoop = true;
          }
        });
    },
    joinSession() {
      this.getSignal(this.inputHash, (resp) => {
        console.log("DOING SIGNAL");
        this.peer.signal(JSON.parse(resp.data));
      });
    },
    getSignal(hash, cb) {
      axios.get('http://localhost:3000/signal/' + hash)
        .then((resp) => {
          if(resp == "bad id") {
            console.error("Bad ID");
            return;
          }

          cb(resp);
        });
    }
  }
}

function initPeer(self) {
  let peer = new Peer({ initiator: location.hash === '#create', trickle: false })
  
  peer.on('error', function (err) { console.log('error', err) })

  peer.on('signal', function (data) {
    console.log("SIGNAL");

    if(self.inputHash === "") {
      self.inputHash = nanoid(10);;
    }

    if(data.type == "offer") {
      self.startAnswerCheck();
      axios.post('http://localhost:3000/signal', {
          id: self.inputHash,
          signal: JSON.stringify(data)
      });
    }else{
      axios.post('http://localhost:3000/answer', {
          id: self.inputHash,
          signal: JSON.stringify(data)
      });
    }
  })

  peer.on('connect', function () {
    console.log("CONNECTED");
    peer.send("New peer connected!");
  })

  peer.on('data', function (data) {
    self.log.push(data.toString());
  })


  return peer;
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
      #outgoing {
        width: 600px;
        word-wrap: break-word;
      }
</style>
