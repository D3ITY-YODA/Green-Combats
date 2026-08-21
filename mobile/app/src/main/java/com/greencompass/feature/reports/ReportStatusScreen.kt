package com.greencompass.feature.reports

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ReportStatusScreen(
    onBack: () -> Unit
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
            Text(text = "Report status", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.md))
            
            Text(text = "Water has changed", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = "Lower Valley", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xxl))
            
            Text(text = "Received", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = "20 Aug, 09:30", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xxl))
            
            Text(text = "Current status", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
            
            TimelineStep(title = "Received", isActive = true, isLast = false)
            TimelineStep(title = "Under review", isActive = true, isLast = false)
            TimelineStep(title = "Verified or resolved", isActive = false, isLast = true)
            
            Spacer(modifier = Modifier.height(AppSpacing.lg))
            Text(text = "We will update this page when the review is complete.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
            
            Spacer(modifier = Modifier.weight(1f))
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}

@Composable
private fun TimelineStep(title: String, isActive: Boolean, isLast: Boolean) {
    Row(modifier = Modifier.height(IntrinsicSize.Min)) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Box(
                modifier = Modifier
                    .size(16.dp)
                    .background(if (isActive) GreenCompassColors.ForestGreen else GreenCompassColors.Stone, shape = CircleShape)
            )
            if (!isLast) {
                Box(
                    modifier = Modifier
                        .width(2.dp)
                        .weight(1f)
                        .background(if (isActive) GreenCompassColors.ForestGreen else GreenCompassColors.Stone)
                )
            }
        }
        Spacer(modifier = Modifier.width(AppSpacing.md))
        Text(
            text = title,
            style = GreenCompassTypography.bodyLarge,
            color = if (isActive) GreenCompassColors.Charcoal else GreenCompassColors.MutedText,
            modifier = Modifier.padding(bottom = AppSpacing.lg)
        )
    }
}
