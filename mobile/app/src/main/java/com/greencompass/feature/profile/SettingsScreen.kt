package com.greencompass.feature.profile

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*
import com.greencompass.navigation.AppRoute

data class SettingsMenuItem(val title: String, val icon: ImageVector)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(
    onBack: () -> Unit,
    onNavigate: (AppRoute) -> Unit,
    onSignOut: () -> Unit
) {
    val menuItems = listOf(
        SettingsMenuItem("Notifications", Icons.Outlined.Notifications),
        SettingsMenuItem("Language", Icons.Outlined.Language),
        SettingsMenuItem("Accessibility", Icons.Outlined.Accessibility),
        SettingsMenuItem("Privacy", Icons.Outlined.PrivacyTip),
        SettingsMenuItem("Data use", Icons.Outlined.DataUsage),
        SettingsMenuItem("About Green Compass", Icons.Outlined.Info)
    )

    GreenCompassScaffold(
        title = "Settings",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            items(menuItems.size) { index ->
                val item = menuItems[index]
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = AppSpacing.lg, vertical = AppSpacing.md)
                        .clickable {
                            when (item.title) {
                                "Notifications" -> onNavigate(AppRoute.NotificationSettings)
                                "Language" -> onNavigate(AppRoute.LanguageSettings)
                                "Accessibility" -> onNavigate(AppRoute.AccessibilitySettings)
                                "Privacy" -> onNavigate(AppRoute.PrivacySettings)
                                "About Green Compass" -> onNavigate(AppRoute.About)
                            }
                        },
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(item.icon, contentDescription = null, tint = GreenCompassColors.Charcoal, modifier = Modifier.size(24.dp))
                    Spacer(Modifier.width(AppSpacing.md))
                    Text(text = item.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.weight(1f))
                }
                if (index < menuItems.size - 1) {
                    Divider(modifier = Modifier.padding(horizontal = AppSpacing.lg), color = GreenCompassColors.Stone)
                }
            }

            item {
                Spacer(modifier = Modifier.height(AppSpacing.xxl))
                TextButton(
                    onClick = onSignOut,
                    modifier = Modifier.fillMaxWidth().padding(horizontal = AppSpacing.lg)
                ) {
                    Text(text = "Sign out", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.EmergencyRed)
                }
                Spacer(modifier = Modifier.height(AppSpacing.xxl))
            }
        }
    }
}
