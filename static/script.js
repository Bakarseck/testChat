const loginButton = document.getElementById('loginButton');

loginButton.addEventListener("click", (event) => {
    event.preventDefault();

    login();
})

const registerButton = document.getElementById('registerButton');

registerButton.addEventListener("click", (event) => {
    event.preventDefault();

    signup();
})


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
                alert("Redirecting to Home page")
            } else {
                alert('Echec Try again !!!');
            }
        });
}

function signup() {
    const nickname = document.querySelector('[name="nickname"]').value;
    const firstName = document.querySelector('[name="firstName"]').value;
    const lastName = document.querySelector('[name="lastName"]').value;
    const age = document.querySelector('[name="age"]').value;
    const gender = document.querySelector('[name="gender"]').value;
    const email = document.querySelector('[name="email"]').value;
    const password = document.querySelector('[name="password"]').value;

    fetch('/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `password=${encodeURIComponent(password)}&username=${encodeURIComponent(nickname)}&firstName=${encodeURIComponent(firstName)}&lastName=${encodeURIComponent(lastName)}&age=${encodeURIComponent(age)}&gender=${encodeURIComponent(gender)}&email=${encodeURIComponent(email)}`
    })
        .then(response => {
            if (response.ok) {
                alert('Inscription réussie, veuillez maintenant vous connecter.');
            } else {
                alert('Echec de l\'inscription, veuillez réessayer.');
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

socket.onopen = function (e) {
    console.log("Connexion au WebSocket établie");
};

socket.onmessage = function (event) {
    let data = JSON.parse(event.data);

    switch (data.Sender) {
        case 'Serveur':
            console.log(data);
            break;
    }
};

socket.onclose = function (event) {
    if (event.wasClean) {
        console.log(`Connexion fermée proprement, code=${event.code} raison=${event.reason}`);
    } else {
        console.log('Connexion fermée de manière inattendue');
    }
};

document.addEventListener('DOMContentLoaded', () => {
    CheckSession();
})

async function CheckSession() {
    fetch("/verify")
        .then(response => {
            if (!response.ok) {
                console.log("no cookie")
                return false;
            } else {
                console.log("cookie present")
                return true;
            }
        });
}