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
fun AccessibilitySettingsScreen(onBack: () -> Unit) {
    var textSize by remember { mutableStateOf("Standard") }
    var highContrast by remember { mutableStateOf(false) }
    var audioUpdates by remember { mutableStateOf(false) }
    var reduceMotion by remember { mutableStateOf(false) }

    GreenCompassScaffold(
        title = "Accessibility",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)) {
            Spacer(modifier = Modifier.height(AppSpacing.xl))

            Text(text = "Text size", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.sm))
            listOf("Standard", "Large", "Extra large").forEach { size ->
                Row(
                    modifier = Modifier.fillMaxWidth().padding(vertical = AppSpacing.xs),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Text(text = size, style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
                    RadioButton(selected = textSize == size, onClick = { textSize = size }, colors = RadioButtonDefaults.colors(selectedColor = GreenCompassColors.ForestGreen))
                }
            }

            Spacer(modifier = Modifier.height(AppSpacing.xl))
            Divider(color = GreenCompassColors.Stone)
            Spacer(modifier = Modifier.height(AppSpacing.md))

            ToggleRow("High contrast", highContrast) { highContrast = it }
            ToggleRow("Audio updates", audioUpdates) { audioUpdates = it }
            ToggleRow("Reduce motion", reduceMotion) { reduceMotion = it }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}

@Composable
private fun ToggleRow(title: String, checked: Boolean, onCheckedChange: (Boolean) -> Unit) {
    Row(
        modifier = Modifier.fillMaxWidth().padding(vertical = AppSpacing.sm),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.SpaceBetween
    ) {
        Text(text = title, style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
        Switch(checked = checked, onCheckedChange = onCheckedChange, colors = SwitchDefaults.colors(checkedThumbColor = GreenCompassColors.ForestGreen, checkedTrackColor = GreenCompassColors.SoftSage))
    }
    Divider(color = GreenCompassColors.Stone)
}
