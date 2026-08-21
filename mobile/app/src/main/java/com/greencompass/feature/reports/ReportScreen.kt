package com.greencompass.feature.reports

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Landscape
import androidx.compose.material.icons.outlined.QuestionMark
import androidx.compose.material.icons.outlined.Thermostat
import androidx.compose.material.icons.outlined.WaterDrop
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

data class ReportOption(val title: String, val subtitle: String, val icon: ImageVector)

@Composable
fun ReportScreen(onSelectType: (String) -> Unit) {
    val options = listOf(
        ReportOption("Weather", "Rain, heat, wind", Icons.Outlined.Thermostat),
        ReportOption("Water", "Rivers, lakes, wells", Icons.Outlined.WaterDrop),
        ReportOption("Land", "Soil, vegetation", Icons.Outlined.Landscape),
        ReportOption("Other", "Anything else", Icons.Outlined.QuestionMark)
    )

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Text(text = "Report", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
            Text(text = "Share what you are seeing in your area.", fontSize = 15.sp, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = 24.dp))
            
            LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                items(options) { option -> 
                    Card(
                        shape = RoundedCornerShape(12.dp), 
                        colors = CardDefaults.cardColors(containerColor = Color.White), 
                        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
                    ) {
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(16.dp)
                                .clickable { onSelectType(option.title) },
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Icon(option.icon, contentDescription = null, tint = GreenCompassColors.ForestGreen, modifier = Modifier.size(32.dp))
                            Spacer(Modifier.width(16.dp))
                            Column {
                                Text(text = option.title, fontSize = 16.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.Charcoal)
                                Text(text = option.subtitle, fontSize = 13.sp, color = GreenCompassColors.MutedText)
                            }
                        }
                    }
                }
            }
        }
    }
}
