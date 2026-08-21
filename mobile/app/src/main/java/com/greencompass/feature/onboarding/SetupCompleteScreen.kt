package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@Composable
fun SetupCompleteScreen(
    onFinish: () -> Unit
) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(horizontal = AppSpacing.lg),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Spacer(modifier = Modifier.weight(1f))

        Icon(
            imageVector = Icons.Default.CheckCircle,
            contentDescription = null,
            tint = GreenCompassColors.ForestGreen,
            modifier = Modifier.size(88.dp)
        )

        Spacer(modifier = Modifier.height(AppSpacing.md))

        Text(
            text = "You're ready to explore.",
            style = GreenCompassTypography.headlineLarge,
            color = GreenCompassColors.Charcoal,
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(bottom = AppSpacing.md)
        )

        Text(
            text = "Green Compass will show information relevant\nto your selected place.",
            style = GreenCompassTypography.bodyLarge,
            color = GreenCompassColors.MutedText,
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(bottom = AppSpacing.xxl)
        )

        Text(
            text = "Place",
            style = GreenCompassTypography.labelMedium,
            color = GreenCompassColors.MutedText,
            modifier = Modifier.padding(bottom = AppSpacing.xs)
        )

        Text(
            text = "Lower Valley",
            style = GreenCompassTypography.titleLarge,
            color = GreenCompassColors.Charcoal,
            modifier = Modifier.padding(bottom = AppSpacing.xxxl)
        )

        Spacer(modifier = Modifier.weight(1f))

        PrimaryButton(
            text = "Open Green Compass",
            onClick = onFinish
        )

        Spacer(modifier = Modifier.height(AppSpacing.xxl))
    }
}
