package com.greencompass.feature.places

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

data class PlaceItem(val id: String, val name: String, val isPrimary: Boolean)

@Composable
fun SavedPlacesScreen(onBack: () -> Unit) {
    val places = listOf(
        PlaceItem("1", "Lower Valley", true),
        PlaceItem("2", "Riverside", false),
        PlaceItem("3", "Upper Highland", false)
    )

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
                Text(text = "Saved places", fontSize = 20.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(24.dp))
            
            LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                items(places) { place ->
                    Card(
                        shape = RoundedCornerShape(12.dp), 
                        colors = CardDefaults.cardColors(containerColor = Color.White),
                        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
                    ) {
                        Row(
                            modifier = Modifier.fillMaxWidth().padding(16.dp),
                            verticalAlignment = Alignment.CenterVertically,
                            horizontalArrangement = Arrangement.SpaceBetween
                        ) {
                            Column {
                                Text(text = place.name, fontSize = 16.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.Charcoal)
                                if (place.isPrimary) Text(text = "Current place", fontSize = 12.sp, color = GreenCompassColors.MutedText)
                            }
                            if (place.isPrimary) Icon(Icons.Default.Check, contentDescription = "Primary", tint = GreenCompassColors.ForestGreen)
                        }
                    }
                }
                item {
                    OutlinedButton(onClick = {}, modifier = Modifier.fillMaxWidth().height(56.dp), shape = RoundedCornerShape(12.dp)) {
                        Text(text = "+ Add place", color = GreenCompassColors.ForestGreen)
                    }
                }
            }
        }
    }
}
