let Game = new class {
    constructor() {
    }

    start(f) {
        this.client = new Iguagile.Client(`ws://${window.location.host}/ws`);
        let b;
        const query = window.location.search.substring(1);
        const vars = query.split("&");
        for (let i = 0; i < vars.length; i++) {
            const pair = vars[i].split("=");
            if (pair[0] === 'token') {
                b = pair[1];
            }
        }
        this.room = JSON.parse(atob(b));
        console.log(this.room)
        f();
    }
}();
