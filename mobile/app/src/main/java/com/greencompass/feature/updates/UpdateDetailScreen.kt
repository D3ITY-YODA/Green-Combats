package com.greencompass.feature.updates

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun UpdateDetailScreen(
    updateId: String,
    onBack: () -> Unit,
    onAcknowledge: () -> Unit,
    onShare: () -> Unit
) {
    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)
        ) {
            item {
                Text(text = "Flood warning", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Heavy rainfall may cause flooding\nwithin 24 hours.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                Text(text = "Lower Valley", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
                Text(text = "Updated 10:00", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xxs))
                Text(text = "Applies until tomorrow", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                Text(text = "This update is provided for your selected place.\nStay informed through official local channels.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                Text(text = "Issued by", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
                Text(text = "Local weather authority", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xxl))
                
                PrimaryButton(text = "I understand", onClick = onAcknowledge, modifier = Modifier.padding(bottom = AppSpacing.sm))
                SecondaryButton(text = "Share an update", onClick = onShare)
                
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}
