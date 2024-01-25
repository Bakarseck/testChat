import { login, signup } from "./auth";

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

function initiateWebSocket() {
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
}