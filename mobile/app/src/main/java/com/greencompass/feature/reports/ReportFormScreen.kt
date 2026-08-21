package com.greencompass.feature.reports

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ReportFormScreen(
    reportType: String,
    onBack: () -> Unit,
    onSubmit: () -> Unit
) {
    var whenHappened by remember { mutableStateOf("Now") }
    var description by remember { mutableStateOf("") }
    val whenOptions = listOf("Now", "Earlier today", "Another time")

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
                Text(text = "Share an update", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                Text(text = "Place", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
                Row(modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.lg), horizontalArrangement = Arrangement.SpaceBetween) {
                    Text(text = "Lower Valley", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
                    TextLinkButton(text = "Change", onClick = { })
                }
                
                Text(text = "What are you seeing?", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
                Text(text = reportType, style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.lg))
                
                Text(text = "When did this happen?", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Row(horizontalArrangement = Arrangement.spacedBy(AppSpacing.sm), modifier = Modifier.padding(bottom = AppSpacing.lg)) {
                    whenOptions.forEach { option ->
                        FilterChip(
                            selected = whenHappened == option,
                            onClick = { whenHappened = option },
                            label = { Text(option) },
                            colors = FilterChipDefaults.filterChipColors(
                                selectedContainerColor = GreenCompassColors.SoftSage,
                                selectedLabelColor = GreenCompassColors.ForestGreen
                            )
                        )
                    }
                }
                
                Text(text = "Add a photo or recording", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                Row(horizontalArrangement = Arrangement.spacedBy(AppSpacing.sm), modifier = Modifier.padding(bottom = AppSpacing.lg)) {
                    SecondaryButton(text = "Add photo", onClick = { }, modifier = Modifier.weight(1f).height(AppSpacing.huge))
                    SecondaryButton(text = "Add audio", onClick = { }, modifier = Modifier.weight(1f).height(AppSpacing.huge))
                }
                
                Text(text = "Tell us more", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xs))
                Text(text = "Optional", style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.sm))
                OutlinedTextField(
                    value = description,
                    onValueChange = { description = it },
                    modifier = Modifier.fillMaxWidth().height(120.dp).padding(bottom = AppSpacing.md),
                    shape = RoundedCornerShape(12.dp)
                )
                
                Text(text = "Your location will be shared with the relevant\nreview team to help verify this update.", style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xl))
                
                PrimaryButton(text = "Send update", onClick = onSubmit)
                Spacer(modifier = Modifier.height(AppSpacing.xxxl))
            }
        }
    }
}
