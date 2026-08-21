package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import com.greencompass.core.ui.*
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GoogleSignInScreen(
    onBack: () -> Unit,
    onSuccess: () -> Unit,
    onUsePhoneEmail: () -> Unit
) {
    var state by remember { mutableStateOf<GoogleSignInState>(GoogleSignInState.Loading) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        delay(1500) // Simulate network delay
        state = GoogleSignInState.Success // For demo, we auto-succeed. Change to Error to test failure state.
    }

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
                .padding(horizontal = AppSpacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Spacer(modifier = Modifier.weight(1f))

            when (state) {
                is GoogleSignInState.Loading -> {
                    CircularProgressIndicator(color = GreenCompassColors.ForestGreen)
                    Spacer(modifier = Modifier.height(AppSpacing.md))
                    Text(
                        text = "Connecting to Google…",
                        style = GreenCompassTypography.bodyLarge,
                        color = GreenCompassColors.MutedText
                    )
                }
                is GoogleSignInState.Success -> {
                    Text(
                        text = "Welcome, Amina",
                        style = GreenCompassTypography.headlineLarge,
                        color = GreenCompassColors.Charcoal,
                        modifier = Modifier.padding(bottom = AppSpacing.sm)
                    )
                    Text(
                        text = "Your Google account is connected.",
                        style = GreenCompassTypography.bodyLarge,
                        color = GreenCompassColors.MutedText,
                        textAlign = TextAlign.Center,
                        modifier = Modifier.padding(bottom = AppSpacing.xxl)
                    )
                    PrimaryButton(text = "Continue", onClick = onSuccess)
                }
                is GoogleSignInState.Error -> {
                    Text(
                        text = "Google sign-in could not be completed.",
                        style = GreenCompassTypography.headlineMedium,
                        color = GreenCompassColors.Charcoal,
                        textAlign = TextAlign.Center,
                        modifier = Modifier.padding(bottom = AppSpacing.xl)
                    )
                    PrimaryButton(
                        text = "Try again",
                        onClick = { state = GoogleSignInState.Loading },
                        modifier = Modifier.padding(bottom = AppSpacing.md)
                    )
                    TextLinkButton(
                        text = "Use phone or email instead",
                        onClick = onUsePhoneEmail
                    )
                }
            }

            Spacer(modifier = Modifier.weight(1f))
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}

sealed class GoogleSignInState {
    object Loading : GoogleSignInState()
    object Success : GoogleSignInState()
    object Error : GoogleSignInState()
}
