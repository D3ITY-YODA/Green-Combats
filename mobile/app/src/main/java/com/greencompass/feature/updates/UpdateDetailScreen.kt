package com.greencompass.feature.updates

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun UpdateDetailScreen(onBack: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
            Spacer(Modifier.height(16.dp))
            
            Text(text = "Water level rising", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
            Spacer(Modifier.height(8.dp))
            Text(text = "Near Ahero, Lower Valley · 5h ago", fontSize = 14.sp, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(24.dp))
            
            Text(text = "Water level at River Nyando is rising slowly. No immediate risk.", fontSize = 16.sp, color = GreenCompassColors.Charcoal, lineHeight = 24.sp)
            Spacer(Modifier.height(24.dp))
            
            Card(shape = RoundedCornerShape(12.dp), colors = CardDefaults.cardColors(containerColor = GreenCompassColors.SoftSage)) {
                Row(modifier = Modifier.padding(16.dp)) {
                    Text(text = "Reported by", fontSize = 12.sp, color = GreenCompassColors.MutedText)
                    Spacer(Modifier.width(8.dp))
                    Text(text = "Community member", fontSize = 12.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.Charcoal)
                }
            }
            
            Spacer(Modifier.weight(1f))
            
            Button(onClick = {}, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) {
                Text(text = "I have an update too", fontSize = 16.sp)
            }
        }
    }
}
