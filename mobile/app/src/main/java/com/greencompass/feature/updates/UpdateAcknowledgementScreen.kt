package com.greencompass.feature.updates

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun UpdateAcknowledgementScreen(
    onBack: () -> Unit,
    onShare: () -> Unit
) {
    var acknowledged by remember { mutableStateOf(false) }

    GreenCompassScaffold(
        title = "",
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

            if (!acknowledged) {
                Text(text = "Did you receive this update?", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                SecondaryButton(text = "Yes, I understand", onClick = { acknowledged = true }, modifier = Modifier.padding(bottom = AppSpacing.sm))
                SecondaryButton(text = "The information is unclear", onClick = { acknowledged = true }, modifier = Modifier.padding(bottom = AppSpacing.sm))
                SecondaryButton(text = "Share what I am seeing", onClick = onShare, modifier = Modifier.padding(bottom = AppSpacing.sm))
            } else {
                Icon(imageVector = Icons.Default.CheckCircle, contentDescription = null, tint = GreenCompassColors.ForestGreen, modifier = Modifier.size(64.dp))
                Spacer(modifier = Modifier.height(AppSpacing.md))
                Text(text = "Thank you.", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Text(text = "Your response helps us understand\nwhether updates are reaching people.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
            }

            Spacer(modifier = Modifier.weight(1f))
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
