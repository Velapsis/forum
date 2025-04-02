// Loop
setInterval(function() {CheckCredentials();}, 1000)

function CheckCredentials() {

    Username = document.getElementById("username").value
    Email = document.getElementById("email").value
    Password = document.getElementById("password").value

    // USERNAME CHECK
    // Check for invalid characters
    if (Username != "") {
        if (!/^[A-Za-z0-9_]+$/.test(Username)) {
            document.getElementById("usernameHint").textContent = "Invalid username: Only letters, numbers and underscores are allowed."
        }
        // Check for length
        else if (Username.length < 4) {
            document.getElementById("usernameHint").textContent = "Username is too short. Min length: 4"
        }
        else if (Username.length > 16) {
            document.getElementById("usernameHint").textContent = "Username is too long. Max length: 16"
        }
        // Valid username
        else {
            document.getElementById("usernameHint").textContent = ""
        }
    }

    // EMAIL CHECK
    if (Email != "") {
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(Email)) {
            document.getElementById("emailHint").textContent = "Invalid email."
        }
        // Valid email
        else {
            document.getElementById("emailHint").textContent = ""
        }
    }

    // PASSWORD CHECK
    // Check for length
    if (Password != "") {
        if (Password.length < 8) {
            document.getElementById("passwordHint").textContent = "Your password is too short. It should be at least 8 characters long."
        }
        // Check for number
        else if (!/[1-9]/.test(Password)) {
            document.getElementById("passwordHint").textContent = "Your password should contain at least one number."
        }
        // Check for caps letter
        else if (!/[A-Z]/.test(Password)) {
            document.getElementById("passwordHint").textContent = "Your password should contain at least one capital letter."
        } 
        // Valid password
        else {
            document.getElementById("passwordHint").textContent = ""
        }
    }
}
