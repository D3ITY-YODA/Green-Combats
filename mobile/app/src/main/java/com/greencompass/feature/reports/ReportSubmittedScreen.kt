package com.greencompass.feature.reports

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@Composable
fun ReportSubmittedScreen(
    onViewStatus: () -> Unit,
    onReturn: () -> Unit
) {
    Column(
        modifier = Modifier.fillMaxSize().padding(horizontal = AppSpacing.lg),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Spacer(modifier = Modifier.weight(1f))
        
        Icon(imageVector = Icons.Default.CheckCircle, contentDescription = null, tint = GreenCompassColors.ForestGreen, modifier = Modifier.size(64.dp))
        Spacer(modifier = Modifier.height(AppSpacing.md))
        
        Text(text = "Update received", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center, modifier = Modifier.padding(bottom = AppSpacing.sm))
        Text(text = "Thank you for sharing what you are seeing.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center, modifier = Modifier.padding(bottom = AppSpacing.xxs))
        Text(text = "Your update will be reviewed by the relevant local team.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center, modifier = Modifier.padding(bottom = AppSpacing.xxl))
        
        Text(text = "Reference", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
        Text(text = "UPD-000123", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xxxl))
        
        PrimaryButton(text = "View report status", onClick = onViewStatus, modifier = Modifier.padding(bottom = AppSpacing.sm))
        TextLinkButton(text = "Return to Today", onClick = onReturn)
        
        Spacer(modifier = Modifier.weight(1f))
        Spacer(modifier = Modifier.height(AppSpacing.xxl))
    }
}
