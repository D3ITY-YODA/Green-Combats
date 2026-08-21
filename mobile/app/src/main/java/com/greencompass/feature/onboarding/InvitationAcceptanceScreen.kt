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
fun InvitationAcceptanceScreen(
    onBack: () -> Unit,
    onAccept: () -> Unit,
    onDecline: () -> Unit
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
                text = "Organization invitation",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            Text(
                text = "You have been invited to join",
                style = GreenCompassTypography.bodyLarge,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.sm)
            )

            Text(
                text = "Lower Valley Water Authority",
                style = GreenCompassTypography.titleLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            Text(
                text = "Role",
                style = GreenCompassTypography.labelMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )

            Text(
                text = "Water reviewer",
                style = GreenCompassTypography.bodyLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xxxl)
            )

            PrimaryButton(
                text = "Accept invitation",
                onClick = onAccept,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            SecondaryButton(
                text = "Decline",
                onClick = onDecline
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
