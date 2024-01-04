function login() {
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
    })
    .then(response => {
        if (response.ok) {
            document.getElementById('form-container').style.display = 'none';
            document.getElementById('chat-container').style.display = 'block';
        } else {
            alert('Echec Try again !!!');
        }
    });
}

function signup() {
    const username = document.getElementById('signup-username').value;
    const password = document.getElementById('signup-password').value;

    fetch('/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
    })
    .then(response => {
        if (response.ok) {
            alert('Inscription réussie, veuillez maintenant vous connecter.');
        } else {
            alert('Echec Try again !!!')
        }
    });
}

function showSignup() {
    document.getElementById('login-form').classList.add('hidden');
    document.getElementById('signup-form').classList.remove('hidden');
}

function showLogin() {
    document.getElementById('signup-form').classList.add('hidden');
    document.getElementById('login-form').classList.remove('hidden');
}

let socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = function(e) {
    console.log("Connexion au WebSocket établie");
};

socket.onmessage = function(event) {
    let data = JSON.parse(event.data);

    switch(data.Sender) {
        case 'Serveur':
            console.log(data);
            break;
    }
};


socket.onclose = function(event) {
    if (event.wasClean) {
        console.log(`Connexion fermée proprement, code=${event.code} raison=${event.reason}`);
    } else {
        console.log('Connexion fermée de manière inattendue');
    }
};

