package com.greencompass.feature.reports

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ReportTypeSelectionScreen(
    reportType: String,
    onBack: () -> Unit,
    onContinue: () -> Unit,
    onChooseDifferent: () -> Unit
) {
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
            Text(text = "Share an update", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
            
            Text(text = "You are reporting:", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = reportType, style = GreenCompassTypography.titleLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xxxl))
            
            PrimaryButton(text = "Continue", onClick = onContinue, modifier = Modifier.padding(bottom = AppSpacing.md))
            TextLinkButton(text = "Choose a different type", onClick = onChooseDifferent)
            
            Spacer(modifier = Modifier.weight(1f))
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
