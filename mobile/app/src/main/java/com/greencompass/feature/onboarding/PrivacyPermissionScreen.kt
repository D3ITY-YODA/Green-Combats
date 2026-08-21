package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PrivacyPermissionScreen(
    onBack: () -> Unit,
    onAllowLocation: () -> Unit,
    onChooseManually: () -> Unit,
    onReadPrivacy: () -> Unit
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
                .padding(horizontal = AppSpacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Spacer(modifier = Modifier.height(AppSpacing.xxl))

            Text(
                text = "Your privacy matters",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            Text(
                text = "Green Compass uses your selected places\nto show relevant information.",
                style = GreenCompassTypography.bodyLarge,
                color = GreenCompassColors.MutedText,
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            Text(
                text = "Location access helps us identify the area\nyou want to follow. You can choose a place manually instead.",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(bottom = AppSpacing.xxxl)
            )

            PrimaryButton(
                text = "Allow location",
                onClick = onAllowLocation,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            SecondaryButton(
                text = "Choose manually",
                onClick = onChooseManually,
                modifier = Modifier.padding(bottom = AppSpacing.md)
            )

            TextLinkButton(
                text = "Read privacy information",
                onClick = onReadPrivacy
            )

            Spacer(modifier = Modifier.weight(1f))
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
