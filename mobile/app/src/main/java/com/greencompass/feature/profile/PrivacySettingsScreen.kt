package com.greencompass.feature.profile

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PrivacySettingsScreen(
    onBack: () -> Unit,
    onDeleteAccount: () -> Unit
) {
    val items = listOf(
        "Location access" to "Manage how your location is used",
        "Community reports" to "Manage report visibility",
        "Contact information" to "Manage phone and email access"
    )

    GreenCompassScaffold(
        title = "Privacy",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            items(items.size) { index ->
                val (title, subtitle) = items[index]
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { }
                        .padding(horizontal = AppSpacing.lg, vertical = AppSpacing.md)
                ) {
                    Text(text = title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                    Text(text = subtitle, style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
                }
                if (index < items.size - 1) {
                    Divider(modifier = Modifier.padding(horizontal = AppSpacing.lg), color = GreenCompassColors.Stone)
                }
            }

            item {
                Spacer(modifier = Modifier.height(AppSpacing.xxl))
                TextButton(
                    onClick = onDeleteAccount,
                    modifier = Modifier.fillMaxWidth().padding(horizontal = AppSpacing.lg)
                ) {
                    Text(text = "Delete account", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.EmergencyRed)
                }
                Spacer(modifier = Modifier.height(AppSpacing.xxl))
            }
        }
    }
}
