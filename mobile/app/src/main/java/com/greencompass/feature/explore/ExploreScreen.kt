package com.greencompass.feature.explore

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
import com.greencompass.domain.model.TopicSection

@Composable
fun ExploreScreen() {
    val sections = listOf(
        TopicSection(key = "local_outlook", title = "Local outlook", description = "Conditions for the coming days"),
        TopicSection(key = "seasonal_information", title = "Seasonal information", description = "Changes that may affect your area"),
        TopicSection(key = "water_outlook", title = "Water outlook", description = "Information about nearby water conditions"),
        TopicSection(key = "community_updates", title = "Community updates", description = "Information shared by people nearby")
    )

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Text(text = "Explore", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
            Text(text = "Information relevant to your selected place.", fontSize = 15.sp, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = 24.dp))
            
            LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                items(sections) { section -> ExploreCard(section) }
            }
        }
    }
}

@Composable
private fun ExploreCard(section: TopicSection) {
    Card(shape = RoundedCornerShape(16.dp), colors = CardDefaults.cardColors(containerColor = Color.White), border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)) {
        Column(modifier = Modifier.padding(20.dp)) {
            Text(text = section.title, fontSize = 16.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            Spacer(modifier = Modifier.height(8.dp))
            Text(text = section.description, fontSize = 15.sp, color = GreenCompassColors.MutedText)
            Spacer(modifier = Modifier.height(12.dp))
            Text(text = "View", fontSize = 14.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.ForestGreen)
        }
    }
}
