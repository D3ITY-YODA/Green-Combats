package com.greencompass.feature.help

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AboutScreen(onBack: () -> Unit) {
    GreenCompassScaffold(
        title = "About Green Compass",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Spacer(modifier = Modifier.height(AppSpacing.xxl))

            Text(
                text = "Green Compass connects trusted environmental information\nwith the places and communities that matter.",
                style = GreenCompassTypography.bodyLarge,
                color = GreenCompassColors.MutedText,
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            Text(
                text = "Know Your Place. Move with Change.",
                style = GreenCompassTypography.titleMedium,
                color = GreenCompassColors.ForestGreen,
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(bottom = AppSpacing.xxl)
            )

            Text(
                text = "Version 1.0.0",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xxl)
            )

            Divider(color = GreenCompassColors.Stone, modifier = Modifier.padding(bottom = AppSpacing.md))

            listOf("Privacy", "Terms", "Data sources", "Open-source acknowledgements").forEach { item ->
                Text(
                    text = item,
                    style = GreenCompassTypography.titleMedium,
                    color = GreenCompassColors.Charcoal,
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = AppSpacing.md)
                )
                Divider(color = GreenCompassColors.Stone)
            }

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}
