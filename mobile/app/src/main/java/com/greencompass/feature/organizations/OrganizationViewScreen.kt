package com.greencompass.feature.organizations

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun OrganizationViewScreen(
    onBack: () -> Unit,
    onOpenOnWeb: () -> Unit,
    onReturnToPersonal: () -> Unit
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
            Text(text = "Lower Valley Water Authority", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = "Water reviewer", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xxl))

            Text(text = "Pending review", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = "4 community updates", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xl))

            Text(text = "Information status", style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = "1 source delayed", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.DelayedGrey, modifier = Modifier.padding(bottom = AppSpacing.xxl))

            PrimaryButton(text = "Open on web", onClick = onOpenOnWeb, modifier = Modifier.padding(bottom = AppSpacing.sm))
            SecondaryButton(text = "Return to personal account", onClick = onReturnToPersonal)

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}
