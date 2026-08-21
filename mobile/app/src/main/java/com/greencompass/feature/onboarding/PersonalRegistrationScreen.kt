package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PersonalRegistrationScreen(
    onBack: () -> Unit,
    onContinueWithGoogle: () -> Unit,
    onContinue: () -> Unit
) {
    var name by remember { mutableStateOf("") }
    var phone by remember { mutableStateOf("") }
    var email by remember { mutableStateOf("") }

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Create your account",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            OutlinedTextField(
                value = name,
                onValueChange = { name = it },
                label = { Text("Name") },
                placeholder = { Text("Optional", color = GreenCompassColors.MutedText) },
                modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.md),
                shape = RoundedCornerShape(12.dp)
            )

            OutlinedTextField(
                value = phone,
                onValueChange = { phone = it },
                label = { Text("Phone number") },
                modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.md),
                shape = RoundedCornerShape(12.dp)
            )

            OutlinedTextField(
                value = email,
                onValueChange = { email = it },
                label = { Text("Email") },
                placeholder = { Text("Optional", color = GreenCompassColors.MutedText) },
                modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.xl),
                shape = RoundedCornerShape(12.dp)
            )

            SecondaryButton(
                text = "Continue with Google",
                onClick = onContinueWithGoogle,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            PrimaryButton(
                text = "Continue",
                onClick = onContinue,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            Spacer(modifier = Modifier.weight(1f))

            Text(
                text = "By continuing, you agree to the Terms and Privacy Policy.",
                style = GreenCompassTypography.bodySmall,
                color = GreenCompassColors.MutedText,
                textAlign = TextAlign.Center,
                modifier = Modifier.fillMaxWidth()
            )
            
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
