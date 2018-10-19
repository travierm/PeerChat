
module.exports = function() {
    let self = {
        data:{}
    };

    self.get = function(key) {
        return self.data[key];
    }

    self.set = function(key, value) {
        self.data[key] = value;
    }

    return self;
}