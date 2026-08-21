package com.greencompass.feature.places

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
fun PlaceDetailsScreen(
    placeId: String,
    onBack: () -> Unit,
    onSetPrimary: () -> Unit,
    onEditName: () -> Unit,
    onRemove: () -> Unit
) {
    var updatesOn by remember { mutableStateOf(true) }
    var smsOff by remember { mutableStateOf(false) }

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)
        ) {
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
            Text(text = "Lower Valley", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = "Primary place", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.ForestGreen, modifier = Modifier.padding(bottom = AppSpacing.xxl))

            Row(
                modifier = Modifier.fillMaxWidth().padding(vertical = AppSpacing.sm),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Column {
                    Text(text = "Updates", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
                    Text(text = "On", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
                }
                Switch(checked = updatesOn, onCheckedChange = { updatesOn = it }, colors = SwitchDefaults.colors(checkedThumbColor = GreenCompassColors.ForestGreen, checkedTrackColor = GreenCompassColors.SoftSage))
            }
            Divider(color = GreenCompassColors.Stone)

            Row(
                modifier = Modifier.fillMaxWidth().padding(vertical = AppSpacing.sm),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Column {
                    Text(text = "SMS", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
                    Text(text = "Off", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
                }
                Switch(checked = smsOff, onCheckedChange = { smsOff = it }, colors = SwitchDefaults.colors(checkedThumbColor = GreenCompassColors.ForestGreen, checkedTrackColor = GreenCompassColors.SoftSage))
            }
            Divider(color = GreenCompassColors.Stone)

            Spacer(modifier = Modifier.height(AppSpacing.xl))

            SecondaryButton(text = "Set as primary", onClick = onSetPrimary, modifier = Modifier.padding(bottom = AppSpacing.sm))
            SecondaryButton(text = "Edit name", onClick = onEditName, modifier = Modifier.padding(bottom = AppSpacing.sm))
            
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
            
            TextLinkButton(text = "Remove place", onClick = onRemove)

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}
