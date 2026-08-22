package com.greencompass.feature.reports

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ChevronRight
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@Composable
fun ReportScreen(onSelectType: (String) -> Unit) {
    val options = listOf(
        "Water has changed",
        "Flooding is visible",
        "Conditions are unusually dry",
        "Plants or crops are under stress",
        "Something else"
    )

    GreenCompassScaffold(title = "Report") { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(start = AppSpacing.lg, end = AppSpacing.lg, top = AppSpacing.sm, bottom = AppSpacing.lg)
        ) {
            Text(text = "Share what you are seeing\nin your area.", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xl))
            
            LazyColumn(verticalArrangement = Arrangement.spacedBy(AppSpacing.sm)) {
                items(options) { option ->
                    Surface(
                        modifier = Modifier.fillMaxWidth().clickable { onSelectType(option) },
                        shape = RoundedCornerShape(16.dp),
                        color = Color.White,
                        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
                    ) {
                        Row(
                            modifier = Modifier.padding(AppSpacing.lg).fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(text = option, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                            Icon(Icons.Default.ChevronRight, contentDescription = null, tint = GreenCompassColors.MutedText)
                        }
                    }
                }
            }
            Spacer(modifier = Modifier.height(AppSpacing.xxxl))
        }
    }
}
