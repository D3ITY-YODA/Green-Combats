package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NotificationPreferencesScreen(
    onBack: () -> Unit,
    onContinue: () -> Unit
) {
    var appNotifs by remember { mutableStateOf(true) }
    var sms by remember { mutableStateOf(false) }
    var voice by remember { mutableStateOf(false) }

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
                text = "Choose how to receive updates",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )
            Text(
                text = "You can change these preferences later.",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            NotificationRow("App notifications", appNotifs) { appNotifs = it }
            NotificationRow("SMS", sms) { sms = it }
            NotificationRow("Voice updates", voice) { voice = it }

            Spacer(modifier = Modifier.height(AppSpacing.md))

            Text(
                text = "Important official updates may be sent\nthrough more than one channel.",
                style = GreenCompassTypography.bodySmall,
                color = GreenCompassColors.MutedText
            )

            Spacer(modifier = Modifier.weight(1f))

            PrimaryButton(
                text = "Continue",
                onClick = onContinue
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}

@Composable
private fun NotificationRow(
    title: String,
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = AppSpacing.sm),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.SpaceBetween
    ) {
        Text(text = title, style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
        Switch(
            checked = checked,
            onCheckedChange = onCheckedChange,
            colors = SwitchDefaults.colors(
                checkedThumbColor = GreenCompassColors.ForestGreen,
                checkedTrackColor = GreenCompassColors.SoftSage
            )
        )
    }
    Divider(color = GreenCompassColors.Stone)
}
