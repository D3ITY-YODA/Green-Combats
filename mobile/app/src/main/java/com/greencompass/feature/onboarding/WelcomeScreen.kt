package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Explore
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.AppSpacing
import com.greencompass.core.ui.GreenCompassColors
import com.greencompass.core.ui.GreenCompassTypography
import com.greencompass.core.ui.PrimaryButton
import com.greencompass.core.ui.TextLinkButton

@Composable
fun WelcomeScreen(onGetStarted: () -> Unit, onChooseLanguage: () -> Unit) {
    Column(
        modifier = Modifier.fillMaxSize().padding(horizontal = AppSpacing.lg).padding(top = AppSpacing.xxxl),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Spacer(modifier = Modifier.weight(1f))

        // Clean, built-in compass icon (88dp)
        Icon(
            imageVector = Icons.Outlined.Explore,
            contentDescription = "Green Compass Logo",
            tint = GreenCompassColors.ForestGreen,
            modifier = Modifier.size(88.dp)
        )

        Text(
            text = "Green Compass",
            style = GreenCompassTypography.titleLarge,
            color = GreenCompassColors.Charcoal,
            modifier = Modifier.padding(bottom = AppSpacing.sm)
        )

        Text(
            text = "Know Your Place.\nMove with Change.",
            style = GreenCompassTypography.headlineLarge,
            color = GreenCompassColors.Charcoal,
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(bottom = AppSpacing.md)
        )

        Text(
            text = "Clear environmental updates for the places\nthat matter to you.",
            style = GreenCompassTypography.bodyLarge,
            color = GreenCompassColors.MutedText,
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(bottom = AppSpacing.huge)
        )

        Spacer(modifier = Modifier.weight(1f))

        PrimaryButton(text = "Get started", onClick = onGetStarted, modifier = Modifier.padding(bottom = AppSpacing.sm))
        TextLinkButton(text = "Choose language", onClick = onChooseLanguage)

        Spacer(modifier = Modifier.height(AppSpacing.xxl))
    }
}
