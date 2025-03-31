// Booleans
var isUsernameValid
var isEmailValid
var isPasswordValid

// Input fields
var Username = document.getElementById("username").value
var Email = document.getElementById("email").value
var Password = document.getElementById("password").value

// Text hints
var UsernameHint = document.getElementById("usernameHint").textContent
var EmailHint = document.getElementById("emailHint").textContent
var PasswordHint = document.getElementById("passwordHint").textContent

// Loop
setInterval(function() {CheckCredentials();}, 1000)

function CheckCredentials() {

    console.log("Checking..")

    // Username check
    if (!/^[A-Za-z0-9_]+$/.test(Username)) {
        UsernameHint = "Invalid username: Only letters, numbers and underscores are allowed."
        isUsernameValid = false
    } else {
        isUsernameValid = true
    }
    if (!Username.length > 4) {
        UsernameHint = "Username is too short. Min length: 4"
        console.log("Username is too short")
        isUsernameValid = false
    } else {
        isUsernameValid = true
    }
    if (!Username.length < 16) {
        UsernameHint = "Username is too long. Max length: 16"
        isUsernameValid = true
    } else {
        isUsernameValid = false
    }

    // Email check
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(Email)) {
        EmailHint = "Invalid email."
        isEmailValid = false
    } else {
        isEmailValid = true
    }

    // Password check
    if (!Password.length < 8) {
        PasswordHint = "Your password is too short. It should be at least 8 characters long."
        isPasswordValid = false
    } else {
        isPasswordValid = true
    }
    if (!/^[1-9]+$/.test(Password)) {
        PasswordHint = "YOur password should contain at least one number."
        isPasswordValid = false
    } else {
        isPasswordValid = true
    }
    if (!/^[A-Z]+$/.test(Password)) {
        PasswordHint = "Your password should contain at least one capital letter."
        isPasswordValid = false
    } else {
        isPasswordValid = true
    }

    //ResetHints()
}

function ResetHints() {
    if (isUsernameValid && isEmailValid && isPasswordValid) {
        UsernameHint = ""
        EmailHint = ""
        PasswordHint = ""
    }
}