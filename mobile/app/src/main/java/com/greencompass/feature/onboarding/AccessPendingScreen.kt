package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AccessPendingScreen(
    onBack: () -> Unit,
    onViewRequest: () -> Unit,
    onReturnToToday: () -> Unit
) {
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
                text = "Access pending",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            Text(
                text = "Your request to join Lower Valley Water Authority\nis waiting for administrator approval.",
                style = GreenCompassTypography.bodyLarge,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xxxl)
            )

            PrimaryButton(
                text = "View request",
                onClick = onViewRequest,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            TextLinkButton(
                text = "Return to Today",
                onClick = onReturnToToday
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
