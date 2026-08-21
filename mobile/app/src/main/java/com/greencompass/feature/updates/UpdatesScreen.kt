package com.greencompass.feature.updates

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors
import com.greencompass.domain.model.PublicUpdate

@Composable
fun UpdatesScreen() {
    val updates = listOf(
        PublicUpdate(id = "1", title = "Water level rising", message = "Water level at River Nyando is rising slowly. No immediate risk.", placeName = "Lower Valley", updatedAt = "5h ago"),
        PublicUpdate(id = "2", title = "Heavy rainfall last night", message = "Reported in West Kano. Roads may be slippery.", placeName = "Lower Valley", updatedAt = "8h ago")
    )

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Text(text = "Updates", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = 24.dp))
            
            if (updates.isEmpty()) {
                Text(text = "No important updates", fontSize = 16.sp, color = GreenCompassColors.MutedText)
            } else {
                LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                    items(updates) { update -> UpdateCard(update) }
                }
            }
        }
    }
}

@Composable
private fun UpdateCard(update: PublicUpdate) {
    Card(shape = RoundedCornerShape(16.dp), colors = CardDefaults.cardColors(containerColor = Color.White), border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)) {
        Column(modifier = Modifier.padding(20.dp)) {
            Text(text = update.title, fontSize = 16.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            Spacer(modifier = Modifier.height(8.dp))
            Text(text = update.message, fontSize = 15.sp, color = GreenCompassColors.MutedText)
            Spacer(modifier = Modifier.height(12.dp))
            Text(text = "${update.placeName} · ${update.updatedAt}", fontSize = 12.sp, color = GreenCompassColors.MutedText)
        }
    }
}
