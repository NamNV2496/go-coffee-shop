document.addEventListener("DOMContentLoaded", function() {
    checkToken()
    ping()
    loadCurrentUser()
})

function checkToken() {
    console.log("check token!")
    token = sessionStorage.getItem('token');
    if (token == null || token == "") {
        window.location.href = "index.html"
    }
}

function loadCurrentUser() {
    currentUser = sessionStorage.getItem('user');
    userElement = document.getElementById('currerntUser')
    userElement.textContent = "Logout: " + currentUser
}

let token = sessionStorage.getItem('token');

async function ping() {
    try {
        const response = await fetch('http://localhost:8083/api/v1/ping', {
            method: "GET",
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to renew token');
        }
        const responseText = await response.text();
        if (responseText) {
            const data = JSON.parse(responseText);
            if (data.code == 401) {
                console.log("Need to renewed token");
                renewToken()
            } else if (data.code == 200) {
                console.log("token is still good");
            }
        } else {
            console.log("token is still good");
        }
    } catch (error) {
        console.error('Error:', error);
        // handle error
    }
}

async function renewToken() {
    try {
        let refresh_token = sessionStorage.getItem('refresh_token');
        console.log("Bearer "+ refresh_token)
        const response = await fetch('http://localhost:8083/api/v1/renewToken', {
            method:"POST",
            headers: {
                'Authorization': `Bearer ${refresh_token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to renew token');
        }

        const responseText = await response.text();
        
        // Log the response text for debugging purposes
        console.log('Response Text:', responseText);

        // Check if the response is not empty before parsing it as JSON
        if (!responseText) {
            throw new Error('Empty response from server');
        }

        const data = JSON.parse(responseText);

        if (data.code == 200) {

            const body = data;
            const token = body.message.token; // assuming you want the first token from the message array
            console.log('Success:', token);
            sessionStorage.setItem('token', token);
            console.log("Token renewed successfully");
            // location.reload();
        }
    } catch (error) {
        console.error('Error:', error);
        // handle error
    }
}

function logout() {
    sessionStorage.setItem('token', "");
    sessionStorage.setItem('refresh_token', "");
    window.location.href = "index.html"
}
