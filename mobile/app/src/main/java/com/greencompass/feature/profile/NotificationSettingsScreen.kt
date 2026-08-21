package com.greencompass.feature.profile

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
fun NotificationSettingsScreen(onBack: () -> Unit) {
    var appNotifs by remember { mutableStateOf(true) }
    var sms by remember { mutableStateOf(false) }
    var voice by remember { mutableStateOf(false) }
    var important by remember { mutableStateOf(true) }

    GreenCompassScaffold(
        title = "Notifications",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            Spacer(modifier = Modifier.height(AppSpacing.xl))

            NotificationRow("App notifications", "On", appNotifs) { appNotifs = it }
            NotificationRow("SMS", "Off", sms) { sms = it }
            NotificationRow("Voice updates", "Off", voice) { voice = it }

            Spacer(modifier = Modifier.height(AppSpacing.xl))
            Divider(color = GreenCompassColors.Stone)
            Spacer(modifier = Modifier.height(AppSpacing.md))

            Text(text = "Important updates", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.sm))
            NotificationRow("Allow important updates", "On", important) { important = it }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}

@Composable
private fun NotificationRow(title: String, status: String, checked: Boolean, onCheckedChange: (Boolean) -> Unit) {
    Row(
        modifier = Modifier.fillMaxWidth().padding(vertical = AppSpacing.sm),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.SpaceBetween
    ) {
        Column {
            Text(text = title, style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
            Text(text = status, style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
        }
        Switch(checked = checked, onCheckedChange = onCheckedChange, colors = SwitchDefaults.colors(checkedThumbColor = GreenCompassColors.ForestGreen, checkedTrackColor = GreenCompassColors.SoftSage))
    }
    Divider(color = GreenCompassColors.Stone)
}
