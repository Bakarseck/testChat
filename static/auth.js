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

export { login, signup };