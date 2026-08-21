package com.greencompass.feature.auth

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun SignInScreen(onBack: () -> Unit, onSignInSuccess: () -> Unit) {
    var phone by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
            Spacer(Modifier.height(16.dp))
            Text(text = "Welcome back", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
            Text(text = "Sign in to your account.", fontSize = 15.sp, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = 32.dp))

            OutlinedTextField(value = phone, onValueChange = { phone = it }, label = { Text("Phone number") }, modifier = Modifier.fillMaxWidth().padding(bottom = 16.dp), shape = RoundedCornerShape(12.dp))
            OutlinedTextField(value = password, onValueChange = { password = it }, label = { Text("Password") }, visualTransformation = PasswordVisualTransformation(), modifier = Modifier.fillMaxWidth().padding(bottom = 32.dp), shape = RoundedCornerShape(12.dp))

            Button(onClick = onSignInSuccess, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) { Text("Sign in", fontSize = 16.sp) }
            Spacer(Modifier.weight(1f))
            TextButton(onClick = {}, modifier = Modifier.fillMaxWidth()) { Text("Don't have an account? Sign up", color = GreenCompassColors.ForestGreen) }
        }
    }
}

@Composable
fun SignUpScreen(onBack: () -> Unit, onSignUpSuccess: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
            Spacer(Modifier.height(16.dp))
            Text(text = "Create your account", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
            Text(text = "Sign up to save your places and preferences.", fontSize = 15.sp, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = 32.dp))

            OutlinedTextField(value = "", onValueChange = {}, label = { Text("Full name") }, modifier = Modifier.fillMaxWidth().padding(bottom = 16.dp), shape = RoundedCornerShape(12.dp))
            OutlinedTextField(value = "", onValueChange = {}, label = { Text("Phone number") }, modifier = Modifier.fillMaxWidth().padding(bottom = 16.dp), shape = RoundedCornerShape(12.dp))
            OutlinedTextField(value = "", onValueChange = {}, label = { Text("Email address") }, modifier = Modifier.fillMaxWidth().padding(bottom = 16.dp), shape = RoundedCornerShape(12.dp))
            OutlinedTextField(value = "", onValueChange = {}, label = { Text("Create password") }, visualTransformation = PasswordVisualTransformation(), modifier = Modifier.fillMaxWidth().padding(bottom = 32.dp), shape = RoundedCornerShape(12.dp))

            Button(onClick = onSignUpSuccess, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) { Text("Sign up", fontSize = 16.sp) }
        }
    }
}
