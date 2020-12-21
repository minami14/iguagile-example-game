let Game = new class {
    constructor() {
    }

    start(f) {
        let api = new Iguagile.RoomApiClient("/api");
        this.client = new Iguagile.Client(`ws://${window.location.host}/ws`);
        const g = this;
        let callback = function (r) {
            g.room = r;
            console.log(`create: ${g.room}`);
            g.client.onConnect = function () {
                console.log('connected')
            };
            f();
        };
        let req = {
            application_name: 'memory-game',
            version: 'beta-0.0.1',
            password: '',
            max_user: 2,
            information: {},
        };
        api.create(req, callback);
    }
}();
