package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PhoneEmailSignInScreen(
    onBack: () -> Unit,
    onContinue: () -> Unit,
    onCreateAccount: () -> Unit
) {
    var identifier by remember { mutableStateOf("") }

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
            Spacer(modifier = Modifier.height(AppSpacing.xxl))

            Text(
                text = "Welcome back",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            OutlinedTextField(
                value = identifier,
                onValueChange = { identifier = it },
                label = { Text("Phone number or email") },
                modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.xl),
                shape = RoundedCornerShape(12.dp)
            )

            PrimaryButton(
                text = "Continue",
                onClick = onContinue,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            Spacer(modifier = Modifier.weight(1f))

            Text(
                text = "Don't have an account?",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )

            TextLinkButton(
                text = "Create an account",
                onClick = onCreateAccount
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
