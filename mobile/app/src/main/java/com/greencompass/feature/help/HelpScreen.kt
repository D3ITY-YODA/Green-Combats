package com.greencompass.feature.help

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HelpScreen(onBack: () -> Unit, onContactSupport: () -> Unit) {
    val items = listOf(
        "How Green Compass works",
        "Using your places",
        "Understanding updates",
        "Sharing an update",
        "Managing notifications",
        "Privacy and safety"
    )

    GreenCompassScaffold(
        title = "Help and support",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            items(items.size) { index ->
                Text(
                    text = items[index],
                    style = GreenCompassTypography.titleMedium,
                    color = GreenCompassColors.Charcoal,
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { }
                        .padding(horizontal = AppSpacing.lg, vertical = AppSpacing.md)
                )
                if (index < items.size - 1) {
                    Divider(modifier = Modifier.padding(horizontal = AppSpacing.lg), color = GreenCompassColors.Stone)
                }
            }

            item {
                Spacer(modifier = Modifier.height(AppSpacing.xxl))
                Column(modifier = Modifier.padding(horizontal = AppSpacing.lg)) {
                    Text(text = "Need more help?", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
                    Text(text = "Contact your local support team.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.md))
                    PrimaryButton(text = "Contact support", onClick = onContactSupport)
                }
                Spacer(modifier = Modifier.height(AppSpacing.xxl))
            }
        }
    }
}
