package com.greencompass.feature.profile

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowRight
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*
import com.greencompass.navigation.AppRoute

data class ProfileMenuItem(val title: String, val subtitle: String, val icon: ImageVector)

@Composable
fun ProfileScreen(
    onNavigate: (AppRoute) -> Unit
) {
    val menuItems = listOf(
        ProfileMenuItem("Your places", "Manage the places you follow", Icons.Outlined.Place),
        ProfileMenuItem("Updates", "Choose how you receive information", Icons.Outlined.Notifications),
        ProfileMenuItem("Language", "English", Icons.Outlined.Language),
        ProfileMenuItem("Accessibility", "Text size, contrast and audio", Icons.Outlined.Accessibility),
        ProfileMenuItem("Privacy", "Manage your information", Icons.Outlined.PrivacyTip),
        ProfileMenuItem("Organizations", "View your organization access", Icons.Outlined.Business),
        ProfileMenuItem("Help", "Get support", Icons.Outlined.Help)
    )

    GreenCompassScaffold(title = "Profile") { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            item {
                Column(modifier = Modifier.fillMaxWidth().padding(AppSpacing.lg)) {
                    Text(text = "Amina Njeri", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
                    Text(text = "Personal account", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
                }
                Divider(color = GreenCompassColors.Stone)
            }

            items(menuItems.size) { index ->
                val item = menuItems[index]
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = AppSpacing.lg, vertical = AppSpacing.md)
                        .clickable {
                            when (item.title) {
                                "Your places" -> onNavigate(AppRoute.Places)
                                "Updates" -> onNavigate(AppRoute.NotificationSettings)
                                "Language" -> onNavigate(AppRoute.LanguageSettings)
                                "Accessibility" -> onNavigate(AppRoute.AccessibilitySettings)
                                "Privacy" -> onNavigate(AppRoute.PrivacySettings)
                                "Organizations" -> onNavigate(AppRoute.Organizations)
                                "Help" -> onNavigate(AppRoute.Help)
                            }
                        },
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(item.icon, contentDescription = null, tint = GreenCompassColors.Charcoal, modifier = Modifier.size(24.dp))
                    Spacer(Modifier.width(AppSpacing.md))
                    Column(modifier = Modifier.weight(1f)) {
                        Text(text = item.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                        if (item.subtitle.isNotEmpty()) {
                            Text(text = item.subtitle, style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
                        }
                    }
                    Icon(Icons.Default.KeyboardArrowRight, contentDescription = null, tint = GreenCompassColors.MutedText)
                }
                if (index < menuItems.size - 1) {
                    Divider(modifier = Modifier.padding(horizontal = AppSpacing.lg), color = GreenCompassColors.Stone)
                }
            }

            item {
                Spacer(modifier = Modifier.height(AppSpacing.xxl))
            }
        }
    }
}
